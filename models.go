package goshua

import (
	"errors"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type CRUDModel interface {
	//NewInstance(int) interface{}
	Add(interface{}) (int64, error)
	GetById(int) (interface{}, error)
	GetAll(map[string]string, []string, []string, []string, int64, int64) ([]interface{}, error)
	UpdateById(interface{}) error
	Delete(int) error
}

type BaseCRUDModel struct {
	NewInstance func(int) interface{}
}

// Add insert a new  into database and returns
// last inserted Id on success.
func (b BaseCRUDModel) Add(m interface{}) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetById retrieves  by Id. Returns error if
// Id doesn't exist
func (b BaseCRUDModel) GetById(id int) (v interface{}, err error) {
	v = b.NewInstance(id)
	err = orm.NewOrm().Read(v)
	return
}

func (b BaseCRUDModel) GetAll(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(b.NewInstance(0))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	l := reflect.New(reflect.SliceOf(reflect.TypeOf(b.NewInstance(0))))
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(l.Interface(), fields...); err == nil {
		l = l.Elem()

		if len(fields) == 0 {
			for i := 0; i < l.Len(); i++ {
				ml = append(ml, l.Index(i).Interface())
			}
		} else {
			// trim unused fields
			for i := 0; i < l.Len(); i++ {
				m := make(map[string]interface{})
				val := l.Index(i)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// Update updates by Id and returns error if
// the record to be updated doesn't exist
func (b BaseCRUDModel) UpdateById(m interface{}) (err error) {
	_, err = orm.NewOrm().Update(m)
	return
}

// Delete deletes  by Id and returns error if
// the record to be deleted doesn't exist
func (b BaseCRUDModel) Delete(id int) (err error) {
	v := b.NewInstance(id)
	_, err = orm.NewOrm().Delete(v)
	return
}
