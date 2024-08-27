package naming

const (
	annotationPrefix      = "percona.com/"
	annotationPrefixMysql = "mysql.percona.com/"
)

const (
	LabelName      = "app.kubernetes.io/name"
	LabelInstance  = "app.kubernetes.io/instance"
	LabelManagedBy = "app.kubernetes.io/managed-by"
	LabelPartOf    = "app.kubernetes.io/part-of"
	LabelComponent = "app.kubernetes.io/component"
)

const (
	LabelMySQLPrimary = annotationPrefixMysql + "primary"
	// FKS: annotationPrefixMysql fails with a 404 error even when slashes are URL encoded
	LabelMySQLRole = "role"
	LabelExposed   = annotationPrefix + "exposed"
)

const (
	FinalizerDeleteSSL         = annotationPrefix + "delete-ssl"
	FinalizerDeletePodsInOrder = annotationPrefix + "delete-mysql-pods-in-order"

	FinalizerDeleteBackup = annotationPrefix + "delete-backup"
)

type AnnotationKey string

func (s AnnotationKey) String() string {
	return string(s)
}

const (
	AnnotationSecretHash       AnnotationKey = annotationPrefix + "last-applied-secret"
	AnnotationConfigHash       AnnotationKey = annotationPrefix + "configuration-hash"
	AnnotationTLSHash          AnnotationKey = annotationPrefix + "last-applied-tls"
	AnnotationPasswordsUpdated AnnotationKey = annotationPrefix + "passwords-updated"
	AnnotationLastConfigHash   AnnotationKey = annotationPrefix + "last-config-hash"
)
