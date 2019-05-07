// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package cron_fbs

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Job struct {
	_tab flatbuffers.Table
}

func GetRootAsJob(buf []byte, offset flatbuffers.UOffsetT) *Job {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Job{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Job) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Job) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Job) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Job) Alias() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Job) Type() JobType {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt8(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Job) MutateType(n JobType) bool {
	return rcv._tab.MutateInt8Slot(8, n)
}

func (rcv *Job) Periodicity() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Job) LastRunStart() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Job) LastRunStop() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Job) LastError() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Job) LogFilePath() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Job) OptionsType() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Job) MutateOptionsType(n byte) bool {
	return rcv._tab.MutateByteSlot(20, n)
}

func (rcv *Job) Options(obj *flatbuffers.Table) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		rcv._tab.Union(obj, o)
		return true
	}
	return false
}

func JobStart(builder *flatbuffers.Builder) {
	builder.StartObject(10)
}
func JobAddName(builder *flatbuffers.Builder, name flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(name), 0)
}
func JobAddAlias(builder *flatbuffers.Builder, alias flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(alias), 0)
}
func JobAddType(builder *flatbuffers.Builder, type_ int8) {
	builder.PrependInt8Slot(2, type_, 0)
}
func JobAddPeriodicity(builder *flatbuffers.Builder, periodicity flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(periodicity), 0)
}
func JobAddLastRunStart(builder *flatbuffers.Builder, lastRunStart flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(lastRunStart), 0)
}
func JobAddLastRunStop(builder *flatbuffers.Builder, lastRunStop flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(lastRunStop), 0)
}
func JobAddLastError(builder *flatbuffers.Builder, lastError flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(6, flatbuffers.UOffsetT(lastError), 0)
}
func JobAddLogFilePath(builder *flatbuffers.Builder, logFilePath flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(7, flatbuffers.UOffsetT(logFilePath), 0)
}
func JobAddOptionsType(builder *flatbuffers.Builder, optionsType byte) {
	builder.PrependByteSlot(8, optionsType, 0)
}
func JobAddOptions(builder *flatbuffers.Builder, options flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(9, flatbuffers.UOffsetT(options), 0)
}
func JobEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
