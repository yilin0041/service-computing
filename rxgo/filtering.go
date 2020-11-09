package rxgo

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"
)

type filteringOperator struct {
	opFunc func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool)
}

func (fop filteringOperator) op(ctx context.Context, o *Observable) {
	in := o.pred.outflow
	out := o.outflow
	var _out []interface{}
	var wg sync.WaitGroup
	go func() {
		end := false
		flag := make(map[interface{}]bool)

		timeStart := time.Now()
		timeSample := time.Now()
		for x := range in {
			timeFromStart := time.Since(timeStart)
			timeSampleFromStart := time.Since(timeSample)
			timeStart = time.Now()
			if end {
				continue
			}
			if o.ignoreElement {
				continue
			}
			if o.sample > 0 && timeSampleFromStart < o.sample {
				continue
			}
			if o.debounce > time.Duration(0) && timeFromStart < o.debounce {
				continue
			}
			xv := reflect.ValueOf(x)
			// send an error to stream if the flip not accept error
			if e, ok := x.(error); ok && !o.flip_accept_error {
				o.sendToFlow(ctx, e, out)
				continue
			}
			o.mu.Lock()
			_out = append(_out, x)
			o.mu.Unlock()
			if o.elementAt > 0 || o.take != 0 || o.skip != 0 || o.last || (o.distinct && flag[xv.Interface()]) {
				continue
			}
			o.mu.Lock()
			flag[xv.Interface()] = true
			o.mu.Unlock()
			switch threading := o.threading; threading {
			case ThreadingDefault:
				if o.sample > 0 {
					timeSample = timeSample.Add(o.sample)
				}
				if fop.opFunc(ctx, o, xv, out) {
					end = true
				}
			case ThreadingIO:
				fallthrough
			case ThreadingComputing:
				wg.Add(1)
				go func() {
					defer wg.Done()
					if fop.opFunc(ctx, o, xv, out) {
						end = true
					}
				}()
			default:
			}
			if o.first {
				break
			}
		}
		if o.last && len(_out) > 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				xv := reflect.ValueOf(_out[len(_out)-1])
				fop.opFunc(ctx, o, xv, out)
			}()
		}
		if o.take != 0 || o.skip != 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var step int
				if o.takeOrSkip {
					step = o.take
				} else {
					step = o.skip
				}
				var newIn []interface{}
				var err error
				if (o.takeOrSkip && step > 0) || (!o.takeOrSkip && step < 0) {
					if !o.takeOrSkip {
						step = len(_out) + step
					}
					if step >= len(_out) || step <= 0 {
						newIn, err = nil, errors.New("OutOfBound")
					} else {
						newIn, err = _out[:step], nil
					}
				} else if (o.takeOrSkip && step < 0) || (!o.takeOrSkip && step > 0) {
					if o.takeOrSkip {
						step = len(_out) + step
					}
					if step >= len(_out) || step <= 0 {
						newIn, err = nil, errors.New("OutOfBound")
					} else {
						newIn, err = _out[step:], nil
					}
				} else {
					newIn, err = nil, errors.New("OutOfBound")
				}
				if err != nil {
					o.sendToFlow(ctx, err, out)
				} else {
					xv := newIn
					for _, val := range xv {
						fop.opFunc(ctx, o, reflect.ValueOf(val), out)
					}
				}
			}()
		}

		if o.elementAt != 0 {
			if o.elementAt < 0 || o.elementAt > len(_out) {
				o.sendToFlow(ctx, errors.New("OutOfBound"), out)
			} else {
				xv := reflect.ValueOf(_out[o.elementAt-1])
				fop.opFunc(ctx, o, xv, out)
			}
		}

		wg.Wait()
		if (o.last || o.first) && len(_out) == 0 && !o.flip_accept_error {
			o.sendToFlow(ctx, errors.New("InputNotFound"), out)
		}
		o.closeFlow(out)
	}()
}

func (parent *Observable) newFilteringObservable(name string) (o *Observable) {
	//new Observable
	o = newObservable()
	o.Name = name

	//chain Observables
	parent.next = o
	o.pred = parent
	o.root = parent.root

	//set options
	o.buf_len = BufferLen
	return o
}

var filteringTotalOperator = filteringOperator{opFunc: func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
	var params = []reflect.Value{x}
	x = params[0]
	if !end {
		end = o.sendToFlow(ctx, x.Interface(), out)
	}
	return
},
}

// Debounce : only emit an item from an Observable if a particular timespan has passed without it emitting another item
func (parent *Observable) Debounce(_debounce time.Duration) (o *Observable) {
	o = parent.newFilteringObservable("debounce")
	o.first, o.last, o.ignoreElement, o.distinct = false, false, false, false
	o.debounce, o.take, o.skip = _debounce, 0, 0
	o.operator = filteringTotalOperator
	return o
}

// Distinct :suppress duplicate items emitted by an Observable
func (parent *Observable) Distinct() (o *Observable) {
	o = parent.newFilteringObservable("distinct")
	o.ignoreElement, o.first, o.last, o.distinct = false, false, false, true
	o.debounce, o.take, o.skip = 0, 0, 0
	o.operator = filteringTotalOperator
	return o
}

// ElementAt :emit only item n emitted by an Observable
func (parent *Observable) ElementAt(index int) (o *Observable) {
	o = parent.newFilteringObservable("elementAt")
	o.first, o.last, o.distinct, o.ignoreElement, o.takeOrSkip = false, false, false, false, false
	o.debounce, o.skip, o.take, o.elementAt = 0, 0, 0, index
	o.operator = filteringTotalOperator
	return
}

// First :emit only the first item, or the first item that meets a condition, from an Observable
func (parent *Observable) First() (o *Observable) {
	o = parent.newFilteringObservable("first")
	o.first, o.last, o.ignoreElement, o.distinct = true, false, false, false
	o.debounce, o.take, o.skip = 0, 0, 0
	o.operator = filteringTotalOperator
	return o

}

// IgnoreElement :do not emit any items from an Observable but mirror its termination notification
func (parent *Observable) IgnoreElement() (o *Observable) {
	o = parent.newFilteringObservable("ignoreElement")
	o.first, o.last, o.distinct, o.ignoreElement = false, false, false, true
	o.debounce, o.take, o.skip = 0, 0, 0
	o.operator = filteringTotalOperator
	return o
}

// Last :emit only the last item emitted by an Observable
func (parent *Observable) Last() (o *Observable) {
	o = parent.newFilteringObservable("last")
	o.first, o.last, o.distinct, o.ignoreElement = false, true, false, false
	o.debounce, o.take, o.skip = 0, 0, 0
	o.operator = filteringTotalOperator
	return o
}

// Sample :emit the most recent item emitted by an Observable within periodic time intervals
func (parent *Observable) Sample(_sample time.Duration) (o *Observable) {
	o = parent.newFilteringObservable("sample")
	o.first, o.last, o.distinct, o.ignoreElement, o.takeOrSkip = false, false, false, false, false
	o.debounce, o.skip, o.take, o.elementAt, o.sample = 0, 0, 0, 0, _sample
	o.operator = filteringTotalOperator
	return o
}

// Skip :suppress the first n items emitted by an Observable
func (parent *Observable) Skip(num int) (o *Observable) {
	o = parent.newFilteringObservable("skip")
	o.first, o.last, o.distinct, o.ignoreElement, o.takeOrSkip = false, false, false, false, false
	o.debounce, o.take, o.skip = 0, 0, num
	o.operator = filteringTotalOperator
	return o
}

// SkipLast :suppress the last n items emitted by an Observable
func (parent *Observable) SkipLast(num int) (o *Observable) {
	o = parent.newFilteringObservable("skipLast")
	o.first, o.last, o.distinct, o.ignoreElement, o.takeOrSkip = false, false, false, false, false
	o.debounce, o.take, o.skip = 0, 0, -num
	o.operator = filteringTotalOperator
	return o
}

// Take :emit only the first n items emitted by an Observable
func (parent *Observable) Take(num int) (o *Observable) {
	o = parent.newFilteringObservable("Take")
	o.first, o.last, o.distinct, o.ignoreElement, o.takeOrSkip = false, false, false, false, true
	o.debounce, o.skip, o.take = 0, 0, num
	o.operator = filteringTotalOperator
	return o
}

// TakeLast :emit only the last n items emitted by an Observable
func (parent *Observable) TakeLast(num int) (o *Observable) {
	o = parent.newFilteringObservable("takeLast")
	o.first, o.last, o.distinct, o.ignoreElement, o.takeOrSkip = false, false, false, false, true
	o.debounce, o.skip, o.take = 0, 0, -num
	o.operator = filteringTotalOperator
	return o
}
