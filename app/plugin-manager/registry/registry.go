package registry

import "github.com/Aarav-S2005/mini-api-gateway/app/plugin-manager/plugin"

type PluginRegistry struct {
	plugins map[string]plugin.Plugin
}

func NewRegistry() *PluginRegistry {
	reg := &PluginRegistry{
		plugins: make(map[string]plugin.Plugin),
	}
	reg.Register(&plugin.CorsPlugin{})
	reg.Register(&plugin.IpFilterPlugin{})
	reg.Register(&plugin.RateLimitPlugin{})
	reg.Register(&plugin.JwtPlugin{})

	return reg
}

func (r *PluginRegistry) Register(p plugin.Plugin) {
	r.plugins[p.Name()] = p
}

func (r *PluginRegistry) Get(name string) (plugin.Plugin, bool) {
	p, ok := r.plugins[name]
	return p, ok
}
