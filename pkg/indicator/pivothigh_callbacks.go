// Code generated by "callbackgen -type PivotHigh"; DO NOT EDIT.

package indicator

import ()

func (inc *PivotHigh) OnUpdate(cb func(value float64)) {
	inc.updateCallbacks = append(inc.updateCallbacks, cb)
}

func (inc *PivotHigh) EmitUpdate(value float64) {
	for _, cb := range inc.updateCallbacks {
		cb(value)
	}
}