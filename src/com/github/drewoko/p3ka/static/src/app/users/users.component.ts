import {Component} from "@angular/core";
import {User, UserService} from "./user.service";

@Component({
    selector: 'users',
    styleUrls: [ './users.component.css' ],
    templateUrl: './users.component.html',
    providers: [
        UserService
    ]
})
export class UsersComponent {

    allUsers: User[] = [];
    peka2TvUsers: User[] = [];
    goodGameUsers: User[] = [];

    constructor(private userService : UserService) {
        userService.getTop(null)
            .subscribe(users => this.allUsers = users)
        userService.getTop('peka2tv')
            .subscribe(users => this.peka2TvUsers = users)
        userService.getTop('goodgame')
            .subscribe(users => this.goodGameUsers = users)
    }
}