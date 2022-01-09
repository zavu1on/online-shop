package schemas

type User struct {
	ID       int
	Login    string
	Password string
}

type Basket struct {
	ID      int
	User_ID int
}

type Good struct {
	ID          int
	Title       string
	Description string
	Price       int
	Count       int
}