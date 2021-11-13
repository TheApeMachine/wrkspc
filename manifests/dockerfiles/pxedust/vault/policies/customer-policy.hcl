policy "customer-internal-policy" {
    rules =<<EOF
    # Allow to read secrets related to Grafana.
    path "customer/internal/grafana" {
        capabilities = ["read"]
    }
    EOF
}