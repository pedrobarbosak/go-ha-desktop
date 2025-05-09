export namespace config {
	
	export class Config {
	    version: number;
	    urls: string[];
	    accessToken: string;
	    pinnedDevices: string[];
	    scanInterval: number;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.urls = source["urls"];
	        this.accessToken = source["accessToken"];
	        this.pinnedDevices = source["pinnedDevices"];
	        this.scanInterval = source["scanInterval"];
	    }
	}

}

export namespace ha {
	
	export class Device {
	    ID: string;
	    Name: string;
	    Type: string;
	    State: boolean;
	    Error: any;
	    Brightness?: number;
	
	    static createFrom(source: any = {}) {
	        return new Device(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.State = source["State"];
	        this.Error = source["Error"];
	        this.Brightness = source["Brightness"];
	    }
	}

}

