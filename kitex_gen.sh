kitex -module github.com/123508/douyinshop -I=./proto bitjump_auth.proto
kitex -module github.com/123508/douyinshop -I=./proto bitjump_cart.proto
kitex -module github.com/123508/douyinshop -I=./proto bitjump_checkout.proto
kitex -module github.com/123508/douyinshop -I=./proto bitjump_order.proto
kitex -module github.com/123508/douyinshop -I=./proto bitjump_payment.proto
kitex -module github.com/123508/douyinshop -I=./proto bitjump_product.proto
kitex -module github.com/123508/douyinshop -I=./proto bitjump_user.proto

mkdir apps
cd apps

mkdir auth
cd auth
kitex -module github.com/123508/douyinshop -service bitjump.douyinshop.auth -use github.com/123508/douyinshop/kitex_gen -I ../../proto bitjump_auth.proto
cd ..

mkdir cart
cd cart
kitex -module github.com/123508/douyinshop -service bitjump.douyinshop.cart -use github.com/123508/douyinshop/kitex_gen -I ../../proto bitjump_cart.proto
cd ..

mkdir checkout
cd checkout
kitex -module github.com/123508/douyinshop -service bitjump.douyinshop.checkout -use github.com/123508/douyinshop/kitex_gen -I ../../proto bitjump_checkout.proto
cd ..

mkdir order
cd order
kitex -module github.com/123508/douyinshop -service bitjump.douyinshop.order -use github.com/123508/douyinshop/kitex_gen -I ../../proto bitjump_order.proto
cd ..

mkdir payment
cd payment
kitex -module github.com/123508/douyinshop -service bitjump.douyinshop.payment -use github.com/123508/douyinshop/kitex_gen -I ../../proto bitjump_payment.proto
cd ..

mkdir product
cd product
kitex -module github.com/123508/douyinshop -service bitjump.douyinshop.product -use github.com/123508/douyinshop/kitex_gen -I ../../proto bitjump_product.proto
cd ..

mkdir user
cd user
kitex -module github.com/123508/douyinshop -service bitjump.douyinshop.user -use github.com/123508/douyinshop/kitex_gen -I ../../proto bitjump_user.proto
cd ..