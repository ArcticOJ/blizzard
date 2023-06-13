package pb

//go:generate protoc --go_out=paths=source_relative:./blizzard/ --go-drpc_out=paths=source_relative:./blizzard ./blizzard.proto
//go:generate protoc --go_out=paths=source_relative:/data/Dev/igloo/igloo/pb/blizzard/ --go-drpc_out=paths=source_relative:/data/Dev/igloo/igloo/pb/blizzard/ ./blizzard.proto
