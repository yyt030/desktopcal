export namespace model {
	
	export class Event {
	    id: string;
	    title: string;
	    description: string;
	    date: string;
	    startTime: string;
	    endTime: string;
	    color: string;
	    completed: boolean;
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Event(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.date = source["date"];
	        this.startTime = source["startTime"];
	        this.endTime = source["endTime"];
	        this.color = source["color"];
	        this.completed = source["completed"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class LunarDay {
	    year: number;
	    month: number;
	    day: number;
	    isLeapMonth: boolean;
	    dayName: string;
	    monthName: string;
	    animal: string;
	    ganZhi: string;
	
	    static createFrom(source: any = {}) {
	        return new LunarDay(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.year = source["year"];
	        this.month = source["month"];
	        this.day = source["day"];
	        this.isLeapMonth = source["isLeapMonth"];
	        this.dayName = source["dayName"];
	        this.monthName = source["monthName"];
	        this.animal = source["animal"];
	        this.ganZhi = source["ganZhi"];
	    }
	}
	export class DayInfo {
	    solarDate: string;
	    lunar?: LunarDay;
	    solarTerm: string;
	    events: Event[];
	    isToday: boolean;
	    isCurrentMonth: boolean;
	    bgColor: string;
	
	    static createFrom(source: any = {}) {
	        return new DayInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.solarDate = source["solarDate"];
	        this.lunar = this.convertValues(source["lunar"], LunarDay);
	        this.solarTerm = source["solarTerm"];
	        this.events = this.convertValues(source["events"], Event);
	        this.isToday = source["isToday"];
	        this.isCurrentMonth = source["isCurrentMonth"];
	        this.bgColor = source["bgColor"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class Settings {
	    opacity: number;
	    width: number;
	    height: number;
	    theme: string;
	    alwaysTop: boolean;
	    fontSize: number;
	    firstDayOfWeek: number;
	    showLunar: boolean;
	    showSolarTerm: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.opacity = source["opacity"];
	        this.width = source["width"];
	        this.height = source["height"];
	        this.theme = source["theme"];
	        this.alwaysTop = source["alwaysTop"];
	        this.fontSize = source["fontSize"];
	        this.firstDayOfWeek = source["firstDayOfWeek"];
	        this.showLunar = source["showLunar"];
	        this.showSolarTerm = source["showSolarTerm"];
	    }
	}

}

