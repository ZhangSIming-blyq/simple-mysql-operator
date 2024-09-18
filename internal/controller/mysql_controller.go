package controller

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/ZhangSIming-blyq/mysql-operator/api/v1"
)

// MySQLReconciler 是 MySQL 控制器的结构体
type MySQLReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile 是控制器的核心逻辑，用于处理 MySQL CR 的状态变化
func (r *MySQLReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	/*
		是的，目前的 Reconcile 逻辑主要涵盖以下两点：

		创建 MySQL 的 Deployment：如果 MySQL CR 对应的 Deployment 不存在，会根据 CR 中定义的 spec 创建一个新的 Deployment。
		控制副本数：如果 MySQL CR 的 spec.size 与现有 Deployment 的副本数不匹配，控制器会更新 Deployment，以确保实际副本数与期望的副本数一致。
		除此之外，Reconcile 函数还会更新 MySQL CR 的 status，将 Deployment 中实际的副本数同步到 MySQL CR 的 ReadyReplicas 字段。
	*/
	log := log.FromContext(ctx)

	// 获取 MySQL 实例: 在 Kubernetes 的控制器中，r.Get 函数只会返回一个具体的对象，而不会返回多个对象。这是因为 req.NamespacedName 包含了特定的 namespace 和 name，代表了一个唯一的资源实例。因此，不会出现 Get 返回多个对象的情况。
	var mysql appsv1alpha1.MySQL
	if err := r.Get(ctx, req.NamespacedName, &mysql); err != nil {
		if errors.IsNotFound(err) {
			// 如果没有找到 MySQL 实例，可能已经被删除，不采取行动
			log.Info("MySQL resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get MySQL")
		return ctrl.Result{}, err
	}

	// 检查 MySQL Secret 是否存在，如果不存在则创建
	var secret corev1.Secret
	secretName := mysql.Name // 使用 MySQL CR 的名称作为 Secret 名称
	err := r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: mysql.Namespace}, &secret)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Secret for MySQL", "Secret.Namespace", mysql.Namespace, "Secret.Name", secretName)

		// 从 MySQL CR 中获取 username 和 password
		username := mysql.Spec.Username
		password := mysql.Spec.Password

		// 创建 Secret
		secret = corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: mysql.Namespace,
			},
			StringData: map[string]string{
				"MYSQL_ROOT_PASSWORD": password, // 使用 CR 中的 password
				"MYSQL_USER":          username, // 使用 CR 中的 username
			},
			Type: corev1.SecretTypeOpaque,
		}

		if err := r.Create(ctx, &secret); err != nil {
			log.Error(err, "Failed to create new Secret", "Secret.Namespace", mysql.Namespace, "Secret.Name", secretName)
			return ctrl.Result{}, err
		}
		// 重新排队 Reconcile
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Secret")
		return ctrl.Result{}, err
	}

	// 检查是否存在 MySQL Deployment，如果不存在则创建
	var deployment appsv1.Deployment
	err = r.Get(ctx, types.NamespacedName{Name: mysql.Name, Namespace: mysql.Namespace}, &deployment)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Deployment for MySQL", "Deployment.Namespace", mysql.Namespace, "Deployment.Name", mysql.Name)
		// 根据CRD的定义创建 Deployment
		dep := r.mysqlDeployment(&mysql)
		if err := r.Create(ctx, dep); err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", mysql.Namespace, "Deployment.Name", mysql.Name)
			return ctrl.Result{}, err
		}
		// 重新排队 Reconcile
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// 确保副本数与期望一致
	size := mysql.Spec.Size
	if *deployment.Spec.Replicas != size {
		deployment.Spec.Replicas = &size
		if err := r.Update(ctx, &deployment); err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", mysql.Namespace, "Deployment.Name", mysql.Name)
			return ctrl.Result{}, err
		}
		// 创建完secret没有必要立刻重建
		return ctrl.Result{}, nil
	}

	// 更新 MySQL 状态
	mysql.Status.ReadyReplicas = deployment.Status.ReadyReplicas
	if err := r.Status().Update(ctx, &mysql); err != nil {
		log.Error(err, "Failed to update MySQL status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// mysqlDeployment 返回定义的 MySQL Deployment
func (r *MySQLReconciler) mysqlDeployment(mysql *appsv1alpha1.MySQL) *appsv1.Deployment {
	labels := map[string]string{"app": "mysql", "mysql_cr": mysql.Name}
	replicas := mysql.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysql.Name,
			Namespace: mysql.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "mysql",
						Image: "mysql:5.7",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 3306,
							Name:          "mysql",
						}},
						Env: []corev1.EnvVar{
							{
								Name: "MYSQL_ROOT_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										Key: "MYSQL_ROOT_PASSWORD",
										LocalObjectReference: corev1.LocalObjectReference{
											Name: mysql.Name,
										},
									},
								},
							},
						},
					}},
				},
			},
		},
	}
	// Set the owner reference for garbage collection
	ctrl.SetControllerReference(mysql, dep, r.Scheme)
	return dep
}

// SetupWithManager 将控制器注册到 Manager 中
func (r *MySQLReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.MySQL{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
