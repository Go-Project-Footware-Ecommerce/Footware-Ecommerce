package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	First_Name      *string            `json:"first_name" validate:"required,min=2,max=30"`
	Last_Name       *string            `json:"last_name"  validate:"required,min=2,max=30"`
	Password        *string            `json:"password"   validate:"required,min=6"`
	Email           *string            `json:"email"      validate:"email,required"`
	Phone           *string            `json:"phone"      validate:"required"`
	Token           *string            `json:"token"`
	Refresh_Token   *string            `josn:"refresh_token"`
	Created_At      time.Time          `json:"created_at"`
	Updated_At      time.Time          `json:"updtaed_at"`
	User_type       *string            `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	User_ID         string             `json:"user_id"`
	UserCart        []ProductUser      `json:"usercart" bson:"usercart"`
	Address_Details []Address          `json:"address" bson:"address"`
	Order_Status    []Order            `json:"orders" bson:"orders"`
}

type Product struct {
	Product_ID   primitive.ObjectID `bson:"_id"`
	Product_Name *string            `json:"product_name"`
	Price        *uint64            `json:"price"`
	Rating       *uint8             `json:"rating"`
	Image        *string            `json:"image"`
	Size         []SizeDetails      `json:"size"`
	Colour       *string            `json:"colour"`
	Brand        *string            `json:"brand"`
	Gender       *string            `json:"gender"`
	Category     *string            `json:"category"`
	Review       []UserReview       `json:"review"`
}

type ProductUser struct {
	Product_ID   primitive.ObjectID `bson:"_id"`
	Product_Name *string            `json:"product_name" bson:"product_name"`
	Price        int                `json:"price"  bson:"price"`
	Rating       *uint              `json:"rating" bson:"rating"`
	Image        *string            `json:"image"  bson:"image"`
	Size         []SizeDetails      `json:"size"`
	Colour       *string            `json:"colour"`
	Brand        *string            `json:"brand"`
	Gender       *string            `json:"gender"`
	Category     *string            `json:"category"`
}

type SizeDetails struct {
	UK55  int `json:"UK5_5" bson:"UK5_5"`
	UK6   int `json:"UK6" bson:"UK6"`
	UK65  int `json:"UK6_5" bson:"UK6_5"`
	UK7   int `json:"UK7" bson:"UK7"`
	UK75  int `json:"UK7_5" bson:"UK7_5"`
	UK8   int `json:"UK8" bson:"UK8"`
	UK85  int `json:"UK8_5" bson:"UK8_5"`
	UK9   int `json:"UK9" bson:"UK9"`
	UK95  int `json:"UK9_5" bson:"UK9_5"`
	UK10  int `json:"UK10" bson:"UK10"`
	UK105 int `json:"UK10_5" bson:"UK10_5"`
	UK11  int `json:"UK11" bson:"UK11"`
	UK115 int `json:"UK11_5" bson:"UK11_5"`
	UK12  int `json:"UK12" bson:"UK12"`
}

type Address struct {
	Address_id  primitive.ObjectID `bson:"_id"`
	House       *string            `json:"house_name" bson:"house_name"`
	Street      *string            `json:"street_name" bson:"street_name"`
	City        *string            `json:"city_name" bson:"city_name"`
	Pincode     *string            `json:"pin_code" bson:"pin_code"`
	AddressType *string            `json:"addresstype"`
}

type Order struct {
	Order_ID       primitive.ObjectID `bson:"_id"`
	Order_Cart     []ProductUser      `json:"order_list"  bson:"order_list"`
	Orderered_At   time.Time          `json:"ordered_on"  bson:"ordered_on"`
	Price          int                `json:"total_price" bson:"total_price"`
	Discount       *int               `json:"discount"    bson:"discount"`
	Payment_Method Payment            `json:"payment_method" bson:"payment_method"`
}

type Payment struct {
	Digital bool `json:"digital" bson:"digital"`
	COD     bool `json:"cod"     bson:"cod"`
}

type UserReview struct {
	UserID   string `json:"userid"`
	Title    string `json:"title"`
	Rating   int    `json:"rating" bson:"rating"`
	Comments string `json:"comments"`
}
