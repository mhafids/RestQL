package rawparser

import (
	"errors"

	"github.com/mhafids/RestQL/repository"

	"github.com/goccy/go-json"
)

// Insert for single raw parser to save Repository ORM
func (mdl *RawModelParser) Insert(data []byte, model interface{}) (repo repository.Repository, err error) {

	if model == nil {
		err = errors.New("Model not Nil")
		return
	}

	mdl.repo.Model(model)

	datamodel := insertmodel{}
	err = json.Unmarshal(data, &datamodel)
	if err != nil {
		return
	}

	modelrepo := repository.Insert{
		Datas: datamodel.Datas,
	}

	mdl.repo.Insert(modelrepo.Datas)

	return
}

// Delete for single raw parser to delete Repository ORM
func (mdl *RawModelParser) Delete(data []byte, model interface{}) (repo repository.Repository, err error) {

	if model == nil {
		err = errors.New("Model not Nil")
		return
	}

	mdl.repo.Model(model)
	datamodel := deletemodel{}
	err = json.Unmarshal(data, &datamodel)
	if err != nil {
		return
	}

	if datamodel.Find.Operator != "" || datamodel.Filter.Operator != "" || datamodel.Where.Operator != "" {
		mdl.filter(datamodel.filtermodel)
	}

	if len(datamodel.Sortby) > 0 || len(datamodel.Sort) > 0 || len(datamodel.Orderby) > 0 {
		mdl.sortby(datamodel.sortbymodel)
	}

	if datamodel.Limit > 0 {
		mdl.limit(datamodel.limitmodel)
	}

	if datamodel.Offset > 0 || datamodel.Skip > 0 {
		mdl.offset(datamodel.offsetmodel)
	}

	repo = mdl.repo

	return
}

// Update for Update data
func (mdl *RawModelParser) Update(data []byte, model interface{}) (repo repository.Repository, err error) {

	if model == nil {
		err = errors.New("Model not Nil")
		return
	}

	mdl.repo.Model(model)

	datamodel := updatemodel{}
	err = json.Unmarshal(data, &datamodel)
	if err != nil {
		return
	}

	modelrepo := repository.Update{
		Datas: datamodel.Datas,
	}

	mdl.repo.Update(modelrepo.Datas)

	if datamodel.Find.Operator != "" || datamodel.Filter.Operator != "" || datamodel.Where.Operator != "" {
		mdl.filter(datamodel.filtermodel)
	}

	if len(datamodel.Sortby) > 0 || len(datamodel.Sort) > 0 || len(datamodel.Orderby) > 0 {
		mdl.sortby(datamodel.sortbymodel)
	}

	if datamodel.Limit > 0 {
		mdl.limit(datamodel.limitmodel)
	}

	repo = mdl.repo

	return
}
