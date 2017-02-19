import {Component, Input} from "@angular/core";
import {Filter, ImageService} from "../../images/image.service";

@Component({
    selector: 'image-filter',
    templateUrl: './image-filter.component.html',
    styleUrls: ['./image-filter.component.css'],
})
export class ImageFilterComponent {

    filter: Filter = 0;

    imageService: ImageService;

    constructor(imageService: ImageService) {
        this.imageService = imageService;
    }

    private setFilter(filter: Filter) {
        this.filter = filter;
        this.imageService.setFilter(filter);
    }
}