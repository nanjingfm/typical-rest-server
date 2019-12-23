package main

// Autogenerated by Typical-Go. DO NOT EDIT.

import (
	"github.com/typical-go/typical-rest-server/app/repository"
	"github.com/typical-go/typical-rest-server/app/service"
	"github.com/typical-go/typical-rest-server/typical"
)

func init() {
	typical.Descriptor.Constructors.Append(repository.NewBookRepo)
	typical.Descriptor.Constructors.Append(repository.NewDataSourceRepo)
	typical.Descriptor.Constructors.Append(repository.NewLocaleRepo)
	typical.Descriptor.Constructors.Append(repository.NewMusicRepo)
	typical.Descriptor.Constructors.Append(repository.NewTagRepo)
	typical.Descriptor.Constructors.Append(service.NewBookService)
	typical.Descriptor.Constructors.Append(service.NewDataSourceService)
	typical.Descriptor.Constructors.Append(service.NewLocaleService)
	typical.Descriptor.Constructors.Append(service.NewMusicService)
	typical.Descriptor.Constructors.Append(service.NewTagService)
}
