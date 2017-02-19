import {Component} from "@angular/core";
import {ImageService, Filter} from "../images/image.service";
import {ImagesComponent} from "../images/images.component";
import {Image} from "../images/image";
import {ImagePageComponent} from "../other/image/image.page.component";
import {Observable} from "rxjs";

@Component({
    selector: 'random',
    templateUrl: 'random.component.html',
    providers: [
        ImageService,
        ImagesComponent
    ]
})
export class RandomImagesComponent extends ImagePageComponent {

    constructor(imageService: ImageService) {
        super(imageService, null);
    }

    protected init() {
        super.load();
    }

    protected requestImages(filter: Filter): Observable<Image[]> {
        return super.getImageService().getRandom(filter);
    }
}