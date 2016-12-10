import {Image} from "../images/image";
import {ImageService} from "../images/image.service";
import {Observable} from "rxjs";
import {ActivatedRoute} from "@angular/router";

export abstract class ImagePageComponent {

    images: Image[];
    imageService: ImageService;
    route: ActivatedRoute;

    constructor(imageService: ImageService, route: ActivatedRoute) {
        this.images = [];
        this.imageService = imageService;
        this.route = route;
        this.init();
    }

    protected getImageService(): ImageService {
        return this.imageService;
    }

    protected abstract init();

    protected addImages(images: Image[]): void {
        this.images = this.images.concat(images);
    }

    protected getRoute(): ActivatedRoute {
        return this.route;
    }

    protected scrollEvent() {
        this.requestImages()
            .subscribe(images => {
                this.addImages(images);
            })
    }

    protected abstract requestImages(): Observable<Image[]>;
}