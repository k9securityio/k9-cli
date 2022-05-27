source = [
  "bin/k9-darwin64",
  "bin/k9-darwinM1"
]
bundle_id = "io.k9security.k9-cli"

apple_id {
  username = "devops@k9security.io"
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "Apple Development: devops@k9security.io (97K55TYCMF)"
}

dmg {
  output_path = "dist/k9-cli.dmg"
  volume_name = "k9-cli"
}

zip {
  output_path = "dist/k9-cli.zip"
}