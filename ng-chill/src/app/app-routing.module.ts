import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { BuildInfoComponent } from './build-info/build-info.component';
import { BuildHistoryComponent } from './build-history/build-history.component';
import { BuildHistoryDetailsComponent } from './build-history-details/build-history-details.component';

const routes: Routes = [
  {
    path: "build/info/list",
    component: BuildInfoComponent,
  },
  {
    path: "build/history/list",
    component: BuildHistoryComponent,
  },
  {
    path: "build/history/details",
    component: BuildHistoryDetailsComponent,
  },
  { path: "**", redirectTo: "/build/info/list" },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { useHash: true })],
  exports: [RouterModule]
})
export class AppRoutingModule { }
