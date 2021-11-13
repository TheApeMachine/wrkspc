# TODO: Kind of questioning if this policy should exist.
# Comes from the Vault documentation as an example.
path "secret/data/*" {
  capabilities = ["create", "update"]
}
