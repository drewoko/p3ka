import {Component, OnDestroy} from "@angular/core";
import {ImagesComponent} from "../images/images.component";
import {ImageService} from "../images/image.service";
import {Image} from "../images/image";
import {ActivatedRoute} from "@angular/router";
import {ImagePageComponent} from "../other/image/image.page.component";
import {Observable, Subscription} from "rxjs";

@Component({
    selector: 'user',
    templateUrl: './user.component.html',
    providers: [
        ImageService,
        ImagesComponent
    ]
})
export class UserComponent extends ImagePageComponent implements OnDestroy {

    id: number;
    user: string;

    routeSubscription: Subscription;

    constructor(imageService: ImageService, route: ActivatedRoute) {
        super(imageService, route);
    }

    protected init() {
        this.routeSubscription = this.getRoute().params.subscribe(params => {
            this.images = [];

            this.id = params['id'];
            this.user = params['user'];

            this.load();
        });
    }

    protected addImages(images: Image[]): void {
        if(images.length > 0) {
            if(this.id != null && this.user == null) {
                this.imageService.forceOpenImage.next(images[0]);
            }
            this.user = images[0].name;
        }
        this.images = this.images.concat(images);
    }

    protected requestImages(): Observable<Image[]> {
        if(this.id != null && this.user == null) {
            return this.getImageService().getByUserImageId(this.images.length, this.id);
        } else {
            return this.getImageService().getByUser(this.images.length, this.user);
        }
    }

    ngOnDestroy(): void {
        this.routeSubscription.unsubscribe()
    }
}