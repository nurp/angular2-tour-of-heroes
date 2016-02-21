import { Injectable } from 'angular2/core';
import { Http} from 'angular2/http';

export interface Hero {
  id: number;
  name: string;
}

@Injectable()
export class HeroService {
  constructor(private _http: Http) {}

  getHeroes() {
    return this._http.get('/api/heroes');
  }
 
 getHero(id: number) {
   return this._http.get('/api/heroes/'+id);
  }
}
