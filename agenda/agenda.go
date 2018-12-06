package agenda

import (
	"github.com/camsiabor/qcom/util"
	"sync"
	"time"
)

type AgendaManager struct {
	lock    sync.RWMutex;
	agendas map[string]*Agenda;
}

type Agenda struct {
	Name string;
	week [8]int;
	slice []interface{};
}

var _agendaManager * AgendaManager = &AgendaManager{
	agendas : make(map[string]*Agenda),
};

func GetAgendaManager() (* AgendaManager) {
	return _agendaManager;
}

func (o * AgendaManager) Init(config map[string]interface{}) {
	o.lock.Lock();
	defer o.lock.Unlock();
	for agendaName, agendaObj := range config {
		var agenda = new(Agenda);
		var agendaMap, ok = agendaObj.(map[string]interface{})
		if (!ok) {
			continue;
		}
		agenda.Name = agendaName;
		o.agendas[agendaName] = agenda;
		var week = agendaMap["week"];
		if (week != nil) {
			var weekarray = week.([]interface{});
			for _, weekone := range weekarray {
				var iweek = util.AsInt(weekone, 0);
				if (iweek < 0 || iweek >= 7) {
					iweek = 0;
				}
				agenda.week[iweek] = 1;
			}
		}
		agenda.slice = util.AsSlice(agendaMap["slice"], 0);
	}
}

func (o * AgendaManager) Get(agendaName string) (agenda * Agenda) {
	o.lock.RLock();
	defer o.lock.RUnlock();
	return o.agendas[agendaName];
}


func (o * Agenda) In(t * time.Time) map[string]interface{} {
	if (t == nil) {
		var n = time.Now();
		t = &n;
	}
	var iweek = int64(t.Weekday());
	if (o.week[iweek] <= 0) {
		return nil;
	}
	var itime = t.Hour() * 100 + t.Minute();
	for _, oneslice := range o.slice {
		var mslice = util.AsMap(oneslice, false);
		if (mslice == nil) {
			continue;
		}
		var start = util.GetInt(mslice, 0, "start");
		var end = util.GetInt(mslice, 0, "end");
		if (itime >= start && itime <= end) {
			return mslice;
		}
	}
	return nil;
}

func (o * Agenda) InGet(t * time.Time, def interface{}, keys ... interface{}) (interface{}) {
	var m = o.In(t);
	if (m == nil) {
		return nil;
	}
	return util.Get(m, def, keys...);
}




