import {Injectable} from "@angular/core";
import {Http, Response, RequestOptions, Headers} from "@angular/http";
import {Observable} from "rxjs/Observable";
import "rxjs/Rx";

@Injectable()
export class UserService {

    constructor(private http: Http) {
    }

    getTop(): Observable<User[]> {
        let options = new RequestOptions({headers: new Headers({'Content-Type': 'application/json'})});

        return this.http.get("/api/top", options)
            .map((resp: Response) => resp.json() as User[])
            .catch(UserService.handleError);
    }

    private static handleError(error: Response | any) {
        return Observable.throw(error.toString());
    }
}

export class User {
    cnt: number;
    name: string;
}