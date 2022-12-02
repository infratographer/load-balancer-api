config_version = 2

name = "$podname"

mode = "single"

dataplaneapi {
  host = "127.0.0.1"
  port = 5555

  user "admin" {
    insecure = true
    password = "adminpwd"
  }

  transaction {
    transaction_dir = "/tmp/haproxy"
  }

  advertised {}
}

haproxy {
  config_file = "/etc/haproxy/haproxy.cfg"
  haproxy_bin = "haproxy"

  reload {
    reload_delay    = 15
    reload_cmd      = "kill SIGUSR 1"
    restart_cmd     = "systemctl restart haproxy"
    reload_strategy = "custom"
  }
}