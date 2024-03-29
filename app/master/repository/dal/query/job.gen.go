// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"github.com/hubogle/Crontab/app/master/repository/dal/model"
)

func newJob(db *gorm.DB) job {
	_job := job{}

	_job.jobDo.UseDB(db)
	_job.jobDo.UseModel(&model.Job{})

	tableName := _job.jobDo.TableName()
	_job.ALL = field.NewField(tableName, "*")
	_job.ID = field.NewInt32(tableName, "id")
	_job.Name = field.NewString(tableName, "name")
	_job.Status = field.NewInt32(tableName, "status")
	_job.Command = field.NewString(tableName, "command")
	_job.CronExpr = field.NewString(tableName, "cronExpr")
	_job.PlanTime = field.NewTime(tableName, "planTime")
	_job.NextTime = field.NewTime(tableName, "nextTime")
	_job.IsDelete = field.NewBool(tableName, "isDelete")
	_job.Created = field.NewTime(tableName, "created")
	_job.Updated = field.NewTime(tableName, "updated")

	_job.fillFieldMap()

	return _job
}

type job struct {
	jobDo jobDo

	ALL      field.Field
	ID       field.Int32
	Name     field.String
	Status   field.Int32
	Command  field.String
	CronExpr field.String
	PlanTime field.Time
	NextTime field.Time
	IsDelete field.Bool
	Created  field.Time
	Updated  field.Time

	fieldMap map[string]field.Expr
}

func (j job) Table(newTableName string) *job {
	j.jobDo.UseTable(newTableName)
	return j.updateTableName(newTableName)
}

func (j job) As(alias string) *job {
	j.jobDo.DO = *(j.jobDo.As(alias).(*gen.DO))
	return j.updateTableName(alias)
}

func (j *job) updateTableName(table string) *job {
	j.ALL = field.NewField(table, "*")
	j.ID = field.NewInt32(table, "id")
	j.Name = field.NewString(table, "name")
	j.Status = field.NewInt32(table, "status")
	j.Command = field.NewString(table, "command")
	j.CronExpr = field.NewString(table, "cronExpr")
	j.PlanTime = field.NewTime(table, "planTime")
	j.NextTime = field.NewTime(table, "nextTime")
	j.IsDelete = field.NewBool(table, "isDelete")
	j.Created = field.NewTime(table, "created")
	j.Updated = field.NewTime(table, "updated")

	j.fillFieldMap()

	return j
}

func (j *job) WithContext(ctx context.Context) *jobDo { return j.jobDo.WithContext(ctx) }

func (j job) TableName() string { return j.jobDo.TableName() }

func (j job) Alias() string { return j.jobDo.Alias() }

func (j *job) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := j.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (j *job) fillFieldMap() {
	j.fieldMap = make(map[string]field.Expr, 10)
	j.fieldMap["id"] = j.ID
	j.fieldMap["name"] = j.Name
	j.fieldMap["status"] = j.Status
	j.fieldMap["command"] = j.Command
	j.fieldMap["cronExpr"] = j.CronExpr
	j.fieldMap["planTime"] = j.PlanTime
	j.fieldMap["nextTime"] = j.NextTime
	j.fieldMap["isDelete"] = j.IsDelete
	j.fieldMap["created"] = j.Created
	j.fieldMap["updated"] = j.Updated
}

func (j job) clone(db *gorm.DB) job {
	j.jobDo.ReplaceDB(db)
	return j
}

type jobDo struct{ gen.DO }

func (j jobDo) Debug() *jobDo {
	return j.withDO(j.DO.Debug())
}

func (j jobDo) WithContext(ctx context.Context) *jobDo {
	return j.withDO(j.DO.WithContext(ctx))
}

func (j jobDo) Clauses(conds ...clause.Expression) *jobDo {
	return j.withDO(j.DO.Clauses(conds...))
}

func (j jobDo) Returning(value interface{}, columns ...string) *jobDo {
	return j.withDO(j.DO.Returning(value, columns...))
}

func (j jobDo) Not(conds ...gen.Condition) *jobDo {
	return j.withDO(j.DO.Not(conds...))
}

func (j jobDo) Or(conds ...gen.Condition) *jobDo {
	return j.withDO(j.DO.Or(conds...))
}

func (j jobDo) Select(conds ...field.Expr) *jobDo {
	return j.withDO(j.DO.Select(conds...))
}

func (j jobDo) Where(conds ...gen.Condition) *jobDo {
	return j.withDO(j.DO.Where(conds...))
}

func (j jobDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *jobDo {
	return j.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (j jobDo) Order(conds ...field.Expr) *jobDo {
	return j.withDO(j.DO.Order(conds...))
}

func (j jobDo) Distinct(cols ...field.Expr) *jobDo {
	return j.withDO(j.DO.Distinct(cols...))
}

func (j jobDo) Omit(cols ...field.Expr) *jobDo {
	return j.withDO(j.DO.Omit(cols...))
}

func (j jobDo) Join(table schema.Tabler, on ...field.Expr) *jobDo {
	return j.withDO(j.DO.Join(table, on...))
}

func (j jobDo) LeftJoin(table schema.Tabler, on ...field.Expr) *jobDo {
	return j.withDO(j.DO.LeftJoin(table, on...))
}

func (j jobDo) RightJoin(table schema.Tabler, on ...field.Expr) *jobDo {
	return j.withDO(j.DO.RightJoin(table, on...))
}

func (j jobDo) Group(cols ...field.Expr) *jobDo {
	return j.withDO(j.DO.Group(cols...))
}

func (j jobDo) Having(conds ...gen.Condition) *jobDo {
	return j.withDO(j.DO.Having(conds...))
}

func (j jobDo) Limit(limit int) *jobDo {
	return j.withDO(j.DO.Limit(limit))
}

func (j jobDo) Offset(offset int) *jobDo {
	return j.withDO(j.DO.Offset(offset))
}

func (j jobDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *jobDo {
	return j.withDO(j.DO.Scopes(funcs...))
}

func (j jobDo) Unscoped() *jobDo {
	return j.withDO(j.DO.Unscoped())
}

func (j jobDo) Create(values ...*model.Job) error {
	if len(values) == 0 {
		return nil
	}
	return j.DO.Create(values)
}

func (j jobDo) CreateInBatches(values []*model.Job, batchSize int) error {
	return j.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (j jobDo) Save(values ...*model.Job) error {
	if len(values) == 0 {
		return nil
	}
	return j.DO.Save(values)
}

func (j jobDo) First() (*model.Job, error) {
	if result, err := j.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Job), nil
	}
}

func (j jobDo) Take() (*model.Job, error) {
	if result, err := j.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Job), nil
	}
}

func (j jobDo) Last() (*model.Job, error) {
	if result, err := j.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Job), nil
	}
}

func (j jobDo) Find() ([]*model.Job, error) {
	result, err := j.DO.Find()
	return result.([]*model.Job), err
}

func (j jobDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Job, err error) {
	buf := make([]*model.Job, 0, batchSize)
	err = j.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (j jobDo) FindInBatches(result *[]*model.Job, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return j.DO.FindInBatches(result, batchSize, fc)
}

func (j jobDo) Attrs(attrs ...field.AssignExpr) *jobDo {
	return j.withDO(j.DO.Attrs(attrs...))
}

func (j jobDo) Assign(attrs ...field.AssignExpr) *jobDo {
	return j.withDO(j.DO.Assign(attrs...))
}

func (j jobDo) Joins(fields ...field.RelationField) *jobDo {
	for _, _f := range fields {
		j = *j.withDO(j.DO.Joins(_f))
	}
	return &j
}

func (j jobDo) Preload(fields ...field.RelationField) *jobDo {
	for _, _f := range fields {
		j = *j.withDO(j.DO.Preload(_f))
	}
	return &j
}

func (j jobDo) FirstOrInit() (*model.Job, error) {
	if result, err := j.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Job), nil
	}
}

func (j jobDo) FirstOrCreate() (*model.Job, error) {
	if result, err := j.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Job), nil
	}
}

func (j jobDo) FindByPage(offset int, limit int) (result []*model.Job, count int64, err error) {
	result, err = j.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = j.Offset(-1).Limit(-1).Count()
	return
}

func (j jobDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = j.Count()
	if err != nil {
		return
	}

	err = j.Offset(offset).Limit(limit).Scan(result)
	return
}

func (j *jobDo) withDO(do gen.Dao) *jobDo {
	j.DO = *do.(*gen.DO)
	return j
}
