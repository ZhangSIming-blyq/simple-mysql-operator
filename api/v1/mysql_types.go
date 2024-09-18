package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MySQLSpec 定义了 MySQL 资源的期望状态
type MySQLSpec struct {
	// MySQL 用户名
	Username string `json:"username"`

	// MySQL 用户密码
	Password string `json:"password"`

	// 数据库名称
	Database string `json:"database"`

	// MySQL 实例的副本数（用于扩展）
	Size int32 `json:"size"`

	// 定期备份的时间表（Cron 表达式）
	BackupSchedule string `json:"backupSchedule"`

	// 备份存储路径
	BackupPath string `json:"backupPath"`
}

// MySQLStatus 定义了 MySQL 资源的当前状态
type MySQLStatus struct {
	// 当前可用的 MySQL 副本数
	ReadyReplicas int32 `json:"readyReplicas"`

	// 最近备份的时间
	LastBackupTime *metav1.Time `json:"lastBackupTime,omitempty"`

	// 资源状态条件（如是否可用、是否需要备份等）
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MySQL 是 MySQL 资源的 Schema
type MySQL struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MySQLSpec   `json:"spec,omitempty"`
	Status MySQLStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MySQLList 包含 MySQL 资源的列表
type MySQLList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MySQL `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MySQL{}, &MySQLList{})
}
