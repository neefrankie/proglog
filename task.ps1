param (
    $stage = 'build'
)

$app_name = 'proglog.exe'
$out_dir = './build'

$exec = "$out_dir/$app_name"

switch ($stage)
{
    'build' { go build -o $exec -v ./cmd/server/main.go }
    'run' { "$exec" | Invoke-Expression }
    'protoc' { protoc api/v1/*.proto --go_out=. --go_opt=paths=source-relative --proto_path=. }
}
