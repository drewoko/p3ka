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

    users : User[];

    constructor(private userService : UserService) {
        userService.getTop()
            .subscribe(users => this.putUsers(users))
    }

    putUsers(users : User[]) {
        this.users = users;
    }
}