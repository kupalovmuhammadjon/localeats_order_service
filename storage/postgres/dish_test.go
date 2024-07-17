package postgres

func newDishRepo() *DishRepo {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}

	return &DishRepo{Db: db}
}

// func TestCreateDish(t *testing.T) {
// 	d := newDishRepo()

// 	dish := pb.ReqCreateDish{
// 		KitchenId:   "413c0067-665a-4a55-b27b-117a188dd5d9",
// 		Name:        "",
// 		Price:       0,
// 		Category:    "",
// 		Ingredients: []string{},
// 		Description: "",
// 		Available:   false,
// 	}

// 	_, err := d.CreateDish(context.Background(), &dish)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestUpdateDish(t *testing.T) {
// 	d := newDishRepo()

// 	dish := pb.ReqUpdateDish{
// 		Id: "c3b42e7e-ab17-4301-ba08-9b246b4d330d",
// 		Name:        "",
// 		Price:       0,
// 		Category:    "",
// 		Ingredients: []string{},
// 		Description: "",
// 		Available:   false,
// 	}

// 	_, err := d.UpdateDish(context.Background(), &dish)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestGetDishes(t *testing.T) {
// 	d := newDishRepo()

// 	filter := pb.Pagination{
// 		Page:  1,
// 		Limit: 10,
// 	}

// 	_, err := d.GetDishes(context.Background(), &filter)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestGetDishById(t *testing.T) {
// 	d := newDishRepo()

// 	id := pb.Id{
// 		Id: "9c41cd2b-fe7c-47be-92e5-4f957963db05",
// 	}

// 	_, err := d.GetDishById(context.Background(), &id)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }
