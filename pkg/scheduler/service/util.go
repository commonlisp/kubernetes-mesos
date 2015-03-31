package service

import (
	"bytes"
	"fmt"
	"hash/crc64"
	"sort"
	"strconv"

	mesos "github.com/mesos/mesos-go/mesosproto"
)

// compute a hashcode for ExecutorInfo that may be used as a reasonable litmus test
// with respect to compatibility across HA schedulers. the intent is that an HA scheduler
// should fail-fast if it doesn't pass this test, rather than generating (potentially many)
// errors at run-time because a Mesos master decides that the ExecutorInfo generated by a
// secondary scheduler doesn't match that of the primary scheduler.
//
// see https://github.com/apache/mesos/blob/0.22.0/src/common/type_utils.cpp#L110
func hashExecutorInfo(info *mesos.ExecutorInfo) uint64 {
	// !!! we specifically do NOT include:
	// - Framework ID because it's a value that's initialized too late for us to use
	// - Executor ID because it's a value that includes a copy of this hash
	buf := &bytes.Buffer{}
	buf.WriteString(info.GetName())
	buf.WriteString(info.GetSource())
	buf.Write(info.Data)

	if info.Command != nil {
		buf.WriteString(info.Command.GetValue())
		buf.WriteString(info.Command.GetUser())
		buf.WriteString(strconv.FormatBool(info.Command.GetShell()))
		if sz := len(info.Command.Arguments); sz > 0 {
			x := make([]string, sz)
			copy(x, info.Command.Arguments)
			sort.Strings(x)
			for _, item := range x {
				buf.WriteString(item)
			}
		}
		if vars := info.Command.Environment.GetVariables(); vars != nil && len(vars) > 0 {
			names := []string{}
			e := make(map[string]string)

			for _, v := range vars {
				if name := v.GetName(); name != "" {
					names = append(names, name)
					e[name] = v.GetValue()
				}
			}
			sort.Strings(names)
			for _, n := range names {
				buf.WriteString(n)
				buf.WriteString("=")
				buf.WriteString(e[n])
			}
		}
		if uris := info.Command.GetUris(); len(uris) > 0 {
			su := []string{}
			for _, uri := range uris {
				su = append(su, fmt.Sprintf("%s%t%t", uri.GetValue(), uri.GetExecutable(), uri.GetExtract()))
			}
			sort.Strings(su)
			for _, uri := range su {
				buf.WriteString(uri)
			}
		}
		//TODO(jdef) add support for Resources and Container
	}
	table := crc64.MakeTable(crc64.ECMA)
	return crc64.Checksum(buf.Bytes(), table)
}