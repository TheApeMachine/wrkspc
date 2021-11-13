policy "developer-backend-policy" {
    rules =<<EOF
    # Allow to read secrets related to Github repositories.
    path "developer/backend/repositories" {
        capabilities = ["read", "write", "update", "list"]
    }
    EOF
}