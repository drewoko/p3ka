import {Injectable} from "@angular/core";
import {Http, Response, URLSearchParams, RequestOptions, Headers} from "@angular/http";
import {Image} from "./image";
import {Observable} from "rxjs/Observable";
import "rxjs/Rx";

@Injectable()
export class ImageService {

    constructor(private http: Http) {
    }

    getLast(start: number): Observable<Image[]> {

        let params: URLSearchParams = new URLSearchParams();
        params.set("start", start.toString());

        let options = new RequestOptions({headers: new Headers({'Content-Type': 'application/json'})});
        options.search = params;

        return this.http.get("/api/last", options)
            .map(ImageService.handleResponse)
            .catch(ImageService.handleError);
    }

    getByUserImageId(start: number, id: number): Observable<Image[]> {

        let params: URLSearchParams = new URLSearchParams();
        params.set("start", start.toString());
        params.set("id", id.toString());

        let options = new RequestOptions({headers: new Headers({'Content-Type': 'application/json'})});
        options.search = params;

        return this.http.get("/api/user/id", options)
            .map(ImageService.handleResponse)
            .catch(ImageService.handleError);
    }

    getByUser(start: number, user: string): Observable<Image[]> {

        let params: URLSearchParams = new URLSearchParams();
        params.set("start", start.toString());
        params.set("user", user);

        let options = new RequestOptions({headers: new Headers({'Content-Type': 'application/json'})});
        options.search = params;

        return this.http.get("/api/user", options)
            .map(ImageService.handleResponse)
            .catch(ImageService.handleError);
    }

    getRandom(): Observable<Image[]> {
        let options = new RequestOptions({headers: new Headers({'Content-Type': 'application/json'})});

        return this.http.get("/api/random", options)
            .map(ImageService.handleResponse)
            .catch(ImageService.handleError);
    }

    private static handleResponse(resp: Response): Image[] {
        let jsonResp = resp.json();
        return jsonResp == null ? [] : jsonResp as Image[];
    }

    private static handleError(error: Response | any) {
        return Observable.throw(error.toString());
    }
}