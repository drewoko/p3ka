import {NgModule} from "@angular/core";
import {BrowserModule} from "@angular/platform-browser";
import {FormsModule} from "@angular/forms";
import {HttpModule} from "@angular/http";
import {RouterModule, Routes} from "@angular/router";
import {AppComponent} from "./general/app.component";
import {MainComponent} from "./main/main.component";
import {InfiniteScrollModule} from "angular2-infinite-scroll";
import {ImagesComponent} from "./images/images.component";
import {UsersComponent} from "./users/users.component";
import {UserComponent} from "./users/user.component";
import {RandomImagesComponent} from "./random/random.component";
import {AboutComponent} from "./about/about.component";
import {LazyLoadImageModule} from 'ng2-lazyload-image';
import {ErrorComponent} from "./other/error.component";

const appRoutes: Routes = [
    {path: '', component: MainComponent},
    {path: 'users', component: UsersComponent},
    {path: 'user/:user', component: UserComponent},
    {path: 'show/:id', component: UserComponent},
    {path: 'random', component: RandomImagesComponent},
    {path: 'about', component: AboutComponent},
    {path: '**', component: ErrorComponent}
];

@NgModule({
    bootstrap: [AppComponent],
    declarations: [
        AppComponent,
        MainComponent,
        ImagesComponent,
        UsersComponent,
        UserComponent,
        RandomImagesComponent,
        AboutComponent,
        ErrorComponent
    ],
    imports: [
        BrowserModule,
        FormsModule,
        HttpModule,
        InfiniteScrollModule,
        LazyLoadImageModule,
        RouterModule.forRoot(appRoutes, {useHash: true})
    ]
})
export class AppModule {}