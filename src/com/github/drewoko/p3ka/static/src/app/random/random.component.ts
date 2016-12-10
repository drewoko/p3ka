import {Component} from "@angular/core";
import {ImageService} from "../images/image.service";
import {ImagesComponent} from "../images/images.component";
import {Image} from "../images/image";
import {ImagePageComponent} from "../other/image.page.component";
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
        super(imageService);
    }

    protected init() {
        super.scrollEvent();
    }

    protected requestImages(): Observable<Image[]> {
        return super.getImageService().getRandom();
    }
}