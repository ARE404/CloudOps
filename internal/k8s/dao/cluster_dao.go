package dao

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ClusterDAO interface {
	ListClusters(ctx context.Context, req *model.ListClustersReq) ([]*model.K8sCluster, int64, error)
	GetClusterByID(ctx context.Context, id int) (*model.K8sCluster, error)
	CreateCluster(ctx context.Context, cluster *model.K8sCluster) error
	UpdateCluster(ctx context.Context, cluster *model.K8sCluster) error
	DeleteCluster(ctx context.Context, id int) error
	ListAllClusters(ctx context.Context) ([]*model.K8sCluster, error)
}

type clusterDAO struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewClusterDAO(db *gorm.DB, logger *zap.Logger) ClusterDAO {
	return &clusterDAO{db: db, logger: logger}
}

func (d *clusterDAO) ListClusters(ctx context.Context, req *model.ListClustersReq) ([]*model.K8sCluster, int64, error) {
	var clusters []*model.K8sCluster
	var total int64
	query := d.db.WithContext(ctx).Model(&model.K8sCluster{}).Where("deleted_at IS NULL")
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&clusters).Error; err != nil {
		return nil, 0, err
	}
	return clusters, total, nil
}

func (d *clusterDAO) GetClusterByID(ctx context.Context, id int) (*model.K8sCluster, error) {
	var cluster model.K8sCluster
	if err := d.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&cluster).Error; err != nil {
		return nil, err
	}
	return &cluster, nil
}

func (d *clusterDAO) CreateCluster(ctx context.Context, cluster *model.K8sCluster) error {
	return d.db.WithContext(ctx).Create(cluster).Error
}

func (d *clusterDAO) UpdateCluster(ctx context.Context, cluster *model.K8sCluster) error {
	return d.db.WithContext(ctx).Save(cluster).Error
}

func (d *clusterDAO) DeleteCluster(ctx context.Context, id int) error {
	return d.db.WithContext(ctx).Model(&model.K8sCluster{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (d *clusterDAO) ListAllClusters(ctx context.Context) ([]*model.K8sCluster, error) {
	var clusters []*model.K8sCluster
	if err := d.db.WithContext(ctx).Where("deleted_at IS NULL AND status = 1").Find(&clusters).Error; err != nil {
		return nil, err
	}
	return clusters, nil
}
