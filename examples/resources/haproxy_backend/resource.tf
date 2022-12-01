resource "haproxy_backend" "my-backend" {
  name = "my-backend"
  balance_algorithm =  "roundrobin"
}