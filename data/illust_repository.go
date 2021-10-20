package data

type IllustRepository struct {
	dao *IllustDao
}

func NewIllustRepository(dao *IllustDao) *IllustRepository {
	return &IllustRepository{dao: dao}
}

func (repo *IllustRepository) IsExists(illustId string) bool {
	return repo.dao.CheckExists(illustId)
}

func (repo *IllustRepository) GetById(illustId string) (*Illust, error) {
	return repo.dao.FindByID(illustId)
}

func (repo *IllustRepository) GetRandom(r18 int, limit int) ([]Illust, error) {
	return repo.dao.Random(r18, limit)
}

func (repo *IllustRepository) Save(illust *Illust) error {
	return repo.dao.Save(illust)
}

// 0: 非限制级
// 1: 限制级
// 2: 混合
func (repo *IllustRepository) Search(r18 int, q string, limit int) ([]Illust, error) {
	return repo.dao.RandomSearch(r18, q, limit)
}
