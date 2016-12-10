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

const appRoutes: Routes = [
    {path: '', component: MainComponent},
    {path: 'users', component: UsersComponent},
    {path: 'user/:user', component: UserComponent},
    {path: 'random', component: RandomImagesComponent},
    {path: 'about', component: AboutComponent}
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
        AboutComponent
    ],
    imports: [
        BrowserModule,
        FormsModule,
        HttpModule,
        InfiniteScrollModule,
        RouterModule.forRoot(appRoutes, {useHash: true})
    ]
})
export class AppModule {
}