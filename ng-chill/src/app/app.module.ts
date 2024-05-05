import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { MaterialModule } from './material/material.module';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { NavComponent } from './nav/nav.component';
import { BuildInfoComponent } from './build-info/build-info.component';
import { BuildHistoryComponent } from './build-history/build-history.component';
import { ControlledPaginatorComponent } from './controlled-paginator/controlled-paginator.component';
import { HttpClientModule } from '@angular/common/http';
import { BuildHistoryDetailsComponent } from './build-history-details/build-history-details.component';

@NgModule({
  declarations: [
    AppComponent,
    NavComponent,
    BuildInfoComponent,
    BuildHistoryComponent,
    ControlledPaginatorComponent,
    BuildHistoryDetailsComponent
  ],
  imports: [
    HttpClientModule,
    BrowserModule,
    AppRoutingModule,
    MaterialModule
  ],
  providers: [
    provideAnimationsAsync()
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
