
# MKP Cinema Ticketing 

This Project is created to make a RESTful API for Online Order Cinema Ticketing 
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
note : this app contain seeder that you don't have to run sql script    
see sql script at 
(https://github.com/zevinza/mkp-cinema-ticketing/blob/ma/script.sql)

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
then you can login with these accounts,
| Username | Password     | Role                |
| :-------- | :------- | :------------------------- |
| `admin@mail.com` | `password` | `admin` |
| `armadamuhammads@mail.com` | `password` | `user`|
| `johndoe@mailcom` | `password` | `user` |

after susccessfully logged in, copy `access_token` from response,  then fill Access Token Header with "Bearer [access_token]"  
note : few endpoints is restricted to user role, to access these endpoints, you must login with administrator account above

## Flowcart User 
![flowcart](https://github.com/zevinza/mkp-cinema-ticketing/blob/ma/flowcart.jpg?raw=true)

- pastikan user telah terregistrasi sebelumnya (atau dapat menggunakan `POST:account/register` untuk mendaftar)
- login menggunakan credential yang sesuai
- user dapat melihat dan memilih jadwal tayang
- setelah menentukan jadwal, user dapat memilih kursi yang diinginkan (dapat lebih dari 1 kursi)
- pada saat user memilih kursi, kursi tersebut akan di hold di cache redis hingga sesi berakhir atau berhasil melakukan pembayaran
- setelah user memeriksa pesanan dengan seksama, user dapat Reservasi di kursi yang telah dipilih
- user melakukan pembayaran sesuai nominal yang tertera
- tiket akan diterbitkan sesuai dengan kursi yang telah dipilih
- perlu diperhatikan bahwa meskipun telah book, ketersediaan kursi tidak dapat benar-benar dipastikan hingga melakukan pembayaran

## Entity Relationship Diagram
this diagram shows table relationship on database   
![ERD](https://github.com/zevinza/mkp-cinema-ticketing/blob/ma/erd.jpg?raw=true)

- Relasi dimulai dari 3 master data yaitu, `Movie`, `CinemaLocation`, dan `Theater`.
- dari ketiga data tersebut dibuatlah `ShowSchedule` dengan menambah atribut lainnya
- setiap `ShowSchedule` akan dibuat `Seat` sesuai dan sejumlah data yang terdapat pada `SeatLayout`, dan ditambah atribut `IsAvailable` yang menandakan kursi dapat dipesan atau tidak
- dari data `ShowSchedule` dan `Seat` yang dipilih `User`, akan dirangkum dalam `Cart` yang memiliki masa tenggang
- `User` tidak dapat Menambah `Cart` sebelum memproses `Cart` yang sudah dipilih (atau dapat dibatalkan)
- Setelah `User` Mengonfirmasi `Cart`, `Transaction` akan dibuat berdasarkan data `Cart` dan `ShowSchedule`
- `Ticket` dibuat berdasarkan data `Transaction` dan `Seat`
- `Transaction.BookingCode` dan `Ticket.IsActivated` akan di set setelah `User` melakukan `TransactionPayment`
- `ReferenceCount` tidak memiliki relasi, table ini digunakan untuk generate ReferenceCount pada beberapa entitas

## Problem, Solution, and Discussion
1. Aplikasi ini dibuat secara sederhana, sehingga ada kemungkinan beberapa permasalahan yang dihadapi mampu diselesaikan dengan bantuan third-party
2. untuk mengakomodir pemesanan tiket agar tidak terjadi Race Condition, solusi yang digunakan pada sistem ini menggunakan cache redis yang memvalidasi agar user lain tidak dapat memilih Seat yang sudah dipilih oleh user sebelumnya, meskipun ada kemungkinan kursi akan dapat kembali dipilih apabila sesi telah berakhir
3. Cache Redis dipilih karena merupakan NO-SQL berbasis key-value yang disimpan di memori, sehingga mengurangi beban database dalam memvalidasi ketersediaan kursi
4. Implementasi Go-Routines pada beberapa fungsi yang berhubungan dengan koneksi ke database, agar pemrosesan data lebih cepat