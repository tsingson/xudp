// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package xudp

type PluginList []Plugin

// Contains returns the index of the plugin in the list.
// Returns -1 of the plugin could not be found.
func (pl PluginList) Index(p Plugin) int {
	for i, plg := range pl {
		if plg == p {
			return i
		}
	}

	return -1
}

// Contains returns true if the given plugin is registered.
func (pl PluginList) Contains(p Plugin) bool { return pl.Index(p) > -1 }

// PayloadSize returns the combined payload size for all plugins.
func (pl PluginList) PayloadSize() int {
	var size int

	for _, plg := range pl {
		size += plg.PayloadSize()
	}

	return size
}

// Clear removes all plugins.
func (pl *PluginList) Clear() { *pl = nil }

// Register registers the given plugin.
func (pl *PluginList) Register(p Plugin) {
	if !pl.Contains(p) {
		*pl = append(*pl, p)
	}
}

// Unregister removes the given plugin from the connection.
func (pl *PluginList) Unregister(p Plugin) {
	idx := pl.Index(p)

	if idx == -1 {
		return
	}

	t := (*pl)[:idx]
	t = append(t, (*pl)[idx+1:]...)
	*pl = t
}
