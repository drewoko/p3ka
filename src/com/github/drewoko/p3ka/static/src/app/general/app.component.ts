import {Component, ViewEncapsulation} from "@angular/core";
import {Router, NavigationEnd} from "@angular/router";

declare var ga: any;

@Component({
    selector: 'app',
    encapsulation: ViewEncapsulation.None,
    styleUrls: ['./app.component.css'],
    templateUrl: './app.component.html'
})
export class AppComponent {
    constructor(public router: Router) {
        router.events.distinctUntilChanged((previous: any, current: any) => {
            if(current instanceof NavigationEnd) {
                return previous.url === current.url;
            }
            return true;
        }).subscribe((x: any) => {
            ga('send', 'pageview', x.url);
        });
    }
}