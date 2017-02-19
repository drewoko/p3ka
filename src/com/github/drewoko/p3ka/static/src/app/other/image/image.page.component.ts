import {Image} from "../../images/image";
import {ImageService, Filter} from "../../images/image.service";
import {Observable} from "rxjs";
import {ActivatedRoute} from "@angular/router";

export abstract class ImagePageComponent {

    images: Image[];
    imageService: ImageService;
    route: ActivatedRoute;
    filter: Filter;

    constructor(imageService: ImageService, route: ActivatedRoute) {
        this.images = [];
        this.imageService = imageService;
        this.route = route;
        this.filter = Filter.ALL;
        this.init();

        this.imageService.imageLoadRequestAnnounced$.subscribe(() => {
            this.load(true)
        });

        this.imageService.filterObsAnnounced$.subscribe(filter => {
            if(this.filter == filter) {
                return;
            }
            this.filter = filter;
            this.images = [];
            this.load();
        })
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
        this.requestImages(this.filter)
            .subscribe(images => {
                if(openNewImage) {
                    this.imageService.openImage(images[0]);
                }
                this.addImages(images);
            })
    }

    protected abstract requestImages(filter: Filter): Observable<Image[]>;
}