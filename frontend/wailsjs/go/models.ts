export namespace game {
	
	export class BestTimes {
	    easy?: number;
	    medium?: number;
	    hard?: number;
	
	    static createFrom(source: any = {}) {
	        return new BestTimes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.easy = source["easy"];
	        this.medium = source["medium"];
	        this.hard = source["hard"];
	    }
	}
	export class Cell {
	    face: number;
	    row: number;
	    col: number;
	    state: string;
	    adjacentMines: number;
	
	    static createFrom(source: any = {}) {
	        return new Cell(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.face = source["face"];
	        this.row = source["row"];
	        this.col = source["col"];
	        this.state = source["state"];
	        this.adjacentMines = source["adjacentMines"];
	    }
	}
	export class Face {
	    cells: Cell[][];
	
	    static createFrom(source: any = {}) {
	        return new Face(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cells = this.convertValues(source["cells"], Cell);
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
	export class BoardState {
	    difficulty: string;
	    boardSize: number;
	    faces: Face[];
	
	    static createFrom(source: any = {}) {
	        return new BoardState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.difficulty = source["difficulty"];
	        this.boardSize = source["boardSize"];
	        this.faces = this.convertValues(source["faces"], Face);
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
	

}

