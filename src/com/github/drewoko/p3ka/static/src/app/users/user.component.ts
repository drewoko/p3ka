import {Component} from "@angular/core";
import {ImagesComponent} from "../images/images.component";
import {ImageService} from "../images/image.service";
import {Image} from "../images/image";
import {ActivatedRoute} from "@angular/router";
import {ImagePageComponent} from "../other/image.page.component";
import {Observable} from "rxjs";

@Component({
    selector: 'user',
    templateUrl: './user.component.html',
    providers: [
        ImageService,
        ImagesComponent
    ]
})
export class UserComponent extends ImagePageComponent {

    user: string;

    constructor(imageService: ImageService, private route: ActivatedRoute) {
        super(imageService);
    }

    protected init() {
        this.route.params.subscribe(params => {
            this.user = params['user'];

            this.scrollEvent();
        });
    }

    protected requestImages(): Observable<Image[]> {
        return this.getImageService().getByUser(this.images.length, this.user);
    }
}