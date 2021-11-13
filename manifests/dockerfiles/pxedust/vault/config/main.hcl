storage "raft" {
    path    = "./vault/data"
    node_id = "node1"
}

storage "file" {
    path = "/tmp/vault-data
}

listener "tcp" {
    address     = "127.0.0.1:8200"
    tls_disable = "false"
    tls_cert_file = "/etc/ssl/private/<your org>-dev.com.fullchain.pem"
    tls_key_file  = "/etc/ssl/private/<your org>.com.privkey.pem"
}

api_addr      = "http://127.0.0.1:8200"
cluster_addr  = "https://127.0.0.1:8201"
ui            = true
disable_mlock = false