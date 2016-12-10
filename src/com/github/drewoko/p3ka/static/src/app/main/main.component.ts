import {Component} from "@angular/core";
import {ImageService} from "../images/image.service";
import {ImagesComponent} from "../images/images.component";
import {Image} from "../images/image";
import {ImagePageComponent} from "../other/image.page.component";
import {Observable} from "rxjs";

@Component({
    selector: 'home',
    styleUrls: ['./main.component.css'],
    templateUrl: './main.component.html',
    providers: [
        ImageService,
        ImagesComponent
    ]
})
export class MainComponent extends ImagePageComponent {

    constructor(imageService: ImageService) {
        super(imageService, null);
    }

    protected init() {
        this.scrollEvent();
    }

    protected requestImages(): Observable<Image[]> {
        return super.getImageService().getLast(this.images.length);
    }
}