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

        this.imageService.imageLoadRequestAnnounced$.subscribe(() => {
            this.load(true)
        });
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

    protected load(openNewImage?: boolean) {
        this.requestImages()
            .subscribe(images => {
                if(openNewImage) {
                    this.imageService.openImage(images[0]);
                }
                this.addImages(images);
            })
    }

    protected abstract requestImages(): Observable<Image[]>;
}