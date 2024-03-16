
# Simple Online Shop 

This Project is created to make a RESTful API for online shop application
written with go programming language and go-fiber framework.

before running this application, make sure you have these Prequisites: 
- Docker Compose
- Database Viewer such as Navicat (optional)

## How to Run
1. Extract `.env`
```sh
cp .env.example .env
```
2. Run Docker Compose
```sh
docker-compose up -d
```
3. Run Go fiber app
```sh
docker-compose exec go go run .
```
4. Visit swagger documentation at 
```sh
localhost:8080
```
## Authorization
There's two type of Authorization, Access Token and Header Token.   
before accessing register or login please fill header token with:
```sh
v0x37KYbJqKodL0393Xa6jXaRTTN2eD1
```
then you can login with this account,
| Username | Password     | Role                |
| :-------- | :------- | :------------------------- |
| `admin@mail.com` | `password` | `admin` |
| `armadamuhammads@mail.com` | `password` | `user`|

after susccessfully logged in, copy access_token from response, then fill Access Token Header with "Bearer [access_token]"  
please note that few endpoints is restricted to user role, to access these endpoints, you must login with administrator account above

## Entity Relationship Diagram
this diagram shows table relationship on database   
![ERD](https://github.com/zevinza/synapsis-ecommerce/blob/main/erd.jpg?raw=true)

- each transaction has many details and payments
- reference_count has no relation, is used to generate invoice_number, reference_number, SKU, etc.
- actually, transaction_detail has no relation with product, it only used to filtering or showing live product
