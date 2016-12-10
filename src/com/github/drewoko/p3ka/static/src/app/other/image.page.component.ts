import {Image} from "../images/image";
import {ImageService} from "../images/image.service";
import {Observable} from "rxjs";

export abstract class ImagePageComponent {

    images: Image[];
    imageService: ImageService;

    constructor(imageService: ImageService) {
        this.images = [];
        this.imageService = imageService;
        this.init();
    }

    protected getImageService(): ImageService {
        return this.imageService;
    }

    protected abstract init();

    protected addImages(images: Image[]): void {
        this.images = this.images.concat(images);
    }

    protected scrollEvent() {
        this.requestImages()
            .subscribe(images => {
                this.addImages(images);
            })
    }

    protected abstract requestImages(): Observable<Image[]>;
}