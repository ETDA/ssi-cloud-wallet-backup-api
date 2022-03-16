import * as Knex from "knex";


export async function up(knex: Knex): Promise<void> {
    return knex.schema.table("vcs", function (table) {
        table.dropColumn("vc")
    })
}

export async function down(knex: Knex): Promise<void> {
    return knex.schema.table("vcs", function (table) {
        table.string('vc', 255).notNullable()
    })
}

