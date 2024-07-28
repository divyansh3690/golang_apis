package data

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDB struct {
	client       *mongo.Client
	mgDBproducts string
}

//an interface allows us better handling of functions and variables defines in the interface.

func GetMongoDBFunctions() *mongoDB {
	return &mongoDB{}
}

func (md *mongoDB) Mongo_Connect() *mongo.Client {
	// MongoDB connection URI
	uri := "mongodb+srv://divyanshnumb:j72yV3mdkHO2MFlk@productsprj.cayjksg.mongodb.net/?retryWrites=true&w=majority&appName=Productsprj"

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Set the client in the struct
	md.client = client
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Client initialized", client)

	fmt.Println("Connected to MongoDB!")
	return client
}

func Open_Collection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("PRODUCTS").Collection(collectionName)

	return collection
}

func (md *mongoDB) Insert_one_mdb(collection string, doc interface{}) error {
	if md.client == nil {
		md.Mongo_Connect()
		defer func() {
			err := md.client.Disconnect(context.TODO())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Connection to MongoDB closed.")
		}()
	}

	md.mgDBproducts = "PRODUCTS"
	coll := md.client.Database(md.mgDBproducts).Collection(collection)
	ctx := context.TODO()
	data, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	fmt.Print(data)
	return nil
}

func (md *mongoDB) Update_one_mdb(collection string, doc interface{}, id int) error {
	if md.client == nil {
		md.Mongo_Connect()
		defer func() {
			err := md.client.Disconnect(context.TODO())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Connection to MongoDB closed.")
		}()
	}

	fmt.Println("Inside update one mongo :", doc)
	md.mgDBproducts = "PRODUCTS"
	_, _, err := md.GetProdByID_mdb(collection, id)
	if err != nil {
		return err
	}

	// ERROR OCCURS HERE AS AFTER UPDATING THE DOC I AM GETTING ITS ID AS 0

	coll := md.client.Database(md.mgDBproducts).Collection(collection)
	ctx := context.TODO()

	updatedProdData := bson.M{"$set": doc}
	fmt.Println("to eb edited data", doc)
	_, err = coll.UpdateOne(ctx, bson.M{"id": id}, updatedProdData)

	if err != nil {
		return err
	}
	return nil
}

func (md *mongoDB) Delete_one_mdb(collection string, id int) error {
	if md.client == nil {
		md.Mongo_Connect()
		defer func() {
			err := md.client.Disconnect(context.TODO())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Connection to MongoDB closed.")
		}()
	}

	coll := md.client.Database(md.mgDBproducts).Collection(collection)

	ctx := context.TODO()
	_, err := coll.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}

// NOTE: GET FUNCTIONS WILL BE DIFFERENT FOR MONGO DB.

func (md *mongoDB) GetAll_mdb(collection string) ([]*Product, error) {
	if md.client == nil {
		md.Mongo_Connect()
		defer func() {
			err := md.client.Disconnect(context.TODO())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Connection to MongoDB closed.")
		}()

	}
	md.mgDBproducts = "PRODUCTS"
	coll := md.client.Database(md.mgDBproducts).Collection(collection)
	ctx := context.TODO()
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var newProd_list = []*Product{}

	var index = 0
	for cursor.Next(ctx) {
		var product bson.M
		index++
		if err = cursor.Decode(&product); err != nil {
			fmt.Print("Error while iterating", err)

		}
		tempProd := Product{
			ID:          int(product["id"].(int32)),
			Name:        product["name"].(string),
			Description: product["description"].(string),
			Price:       float32(product["price"].(float64)),
			SKU:         product["sku"].(string),
			CreatedOn:   product["createdon"].(string),
			UpdatedOn:   product["updatedon"].(string),
			DeletedOn:   product["deletedon"].(string),
		}
		newProd_list = append(newProd_list, &tempProd)
		fmt.Println(index, product)
	}
	fmt.Println(newProd_list)
	return newProd_list, nil
}

func (md *mongoDB) GetProdByID_mdb(collection string, id int) (string, *Product, error) {
	if md.client == nil {
		md.Mongo_Connect()
		defer func() {
			err := md.client.Disconnect(context.TODO())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Connection to MongoDB closed.")
		}()

	}
	md.mgDBproducts = "PRODUCTS"
	coll := md.client.Database(md.mgDBproducts).Collection(collection)
	ctx := context.TODO()
	cursor := coll.FindOne(ctx, bson.M{"id": id})

	var productFiltered bson.M

	if err := cursor.Decode(&productFiltered); err != nil {
		return "", nil, fmt.Errorf("no product found with given ID")
	}

	// in case the product exists we send over the mongo db id ,
	if productFiltered != nil {
		fmt.Println(productFiltered["_id"])
		//converting the _id to string ]
		mongoId := productFiltered["_id"]
		stringObjectID := mongoId.(primitive.ObjectID).Hex()

		tempProd := Product{
			ID:          int(productFiltered["id"].(int32)),
			Name:        productFiltered["name"].(string),
			Description: productFiltered["description"].(string),
			Price:       float32(productFiltered["price"].(float64)),
			SKU:         productFiltered["sku"].(string),
			CreatedOn:   productFiltered["createdon"].(string),
			UpdatedOn:   productFiltered["updatedon"].(string),
			DeletedOn:   productFiltered["deletedon"].(string),
		}
		fmt.Println(tempProd)

		return stringObjectID, &tempProd, nil
	}
	return "", nil, fmt.Errorf("no product found with given ID")

}

func Get_UserByEmailorID(email string, user_id int) (*User, error) {

	if email == "" && user_id <= 0 {
		return nil, fmt.Errorf("both email and userID cant be null")
	}
	var decodedUser User
	if !(user_id <= 0) {

		// search by userID
		ctx := context.TODO()
		// var userDef User
		err := collection.FindOne(ctx, bson.M{"userid": user_id}).Decode(&decodedUser)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("no product found with given ID")
			}
			return nil, fmt.Errorf("failed to find product: %v", err)
		}
		return &decodedUser, nil
	} else {

		// search by email

		ctx := context.TODO()
		err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&decodedUser)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("no product found with given email")
			}
			return nil, fmt.Errorf("failed to find product: %v", err)
		}
		return &decodedUser, nil
	}

}
